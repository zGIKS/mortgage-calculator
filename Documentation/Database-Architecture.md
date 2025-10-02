# Arquitectura de Base de Datos: Bounded Contexts IAM y Mortgage

## Introducción

La arquitectura de base de datos de esta aplicación financiera se organiza siguiendo el patrón de Bounded Contexts del Domain-Driven Design (DDD), donde cada contexto delimitado encapsula un modelo de dominio específico con sus propias reglas de negocio, entidades y servicios (Evans, 2003). En este sistema, se han identificado dos bounded contexts principales: IAM (Identity and Access Management) para la gestión de usuarios, y Mortgage para el cálculo y gestión de hipotecas. Esta separación permite mantener la cohesión interna de cada dominio mientras se reduce el acoplamiento entre ellos, facilitando la evolución independiente de cada contexto (Vernon, 2016). La persistencia se implementa utilizando GORM como ORM para PostgreSQL, con modelos de dominio que se mapean a entidades de base de datos relacional. Los repositorios actúan como intermediarios entre la capa de dominio y la infraestructura de persistencia, implementando operaciones CRUD (Create, Read, Update, Delete) de manera consistente con los principios de la arquitectura hexagonal (Alonso, 2018).

## Bounded Context IAM (Users)

El bounded context IAM se encarga exclusivamente de la gestión de identidades de usuario, incluyendo autenticación y autorización. Este contexto mantiene una separación clara de responsabilidades, enfocándose únicamente en aspectos relacionados con usuarios sin interferir con la lógica de negocio de hipotecas. El modelo de dominio define una entidad User que encapsula propiedades como identificador único, correo electrónico, contraseña hasheada y nombre completo, junto con métodos para validar y manipular estos datos de manera segura.

### Modelo de Persistencia: UserModel

El modelo de persistencia UserModel representa la tabla "users" en la base de datos PostgreSQL y se define en el archivo `internal/iam/infrastructure/persistence/models/user_model.go`. Esta estructura utiliza GORM tags para mapear campos de Go a columnas de base de datos, implementando restricciones de integridad como índices únicos y campos no nulos. El campo ID se define como UUID para garantizar unicidad global y evitar conflictos en sistemas distribuidos, mientras que el campo Email tiene un índice único para optimizar búsquedas y prevenir duplicados. Los campos CreatedAt y UpdatedAt se gestionan automáticamente por GORM, proporcionando trazabilidad temporal de las operaciones. El método TableName() especifica explícitamente el nombre de la tabla como "users", siguiendo convenciones de nomenclatura consistentes. Los métodos ToEntity() y FromEntity() implementan el patrón de mapeo entre el modelo de persistencia y la entidad de dominio, convirtiendo tipos primitivos de base de datos a objetos de valor del dominio como UserID y Email, que incluyen validaciones de negocio.

### Repositorio: Operaciones CRUD de Users

El repositorio de usuarios, implementado en `internal/iam/infrastructure/persistence/repositories/user_repository_impl.go`, sigue el patrón Repository para abstraer el acceso a datos y mantener la independencia de la tecnología de persistencia. La interfaz UserRepository define contratos para operaciones CRUD que el repositorio concreto implementa utilizando GORM. El método Save() maneja tanto la creación como la actualización de usuarios, generando automáticamente UUIDs para nuevos registros y actualizando timestamps apropiadamente. Los métodos FindByID() y FindByEmail() permiten recuperar usuarios por diferentes criterios, retornando nil cuando no se encuentran registros para evitar excepciones y facilitar el manejo de casos de uso como login. El método ExistsByEmail() optimiza verificaciones de unicidad sin necesidad de cargar toda la entidad, mejorando el rendimiento en operaciones de registro. Todas las operaciones aceptan un contexto para soporte de cancelación y timeouts, siguiendo mejores prácticas de Go para aplicaciones concurrentes. La implementación utiliza transacciones implícitas de GORM para garantizar consistencia, aunque las operaciones individuales de usuario no requieren transacciones complejas debido a su naturaleza atómica.

## Bounded Context Mortgage (Mortgages y Payment Schedule Items)

El bounded context Mortgage encapsula toda la lógica relacionada con el cálculo y gestión de hipotecas, incluyendo la generación de cronogramas de pago y métricas financieras avanzadas. Este contexto mantiene independencia del contexto IAM, referenciando usuarios únicamente por ID sin acceder a sus detalles personales. El modelo de dominio incluye entidades como Mortgage y PaymentSchedule, con value objects para tipos como Currency y RateType que garantizan invariantes de negocio.

### Modelos de Persistencia: MortgageModel y PaymentScheduleItemModel

Los modelos de persistencia para el contexto Mortgage se definen en archivos separados para mantener la claridad y facilitar el mantenimiento. El MortgageModel, ubicado en `internal/mortgage/infrastructure/persistence/models/mortgage_model.go`, representa la tabla "mortgages" y contiene todos los campos necesarios para almacenar una hipoteca completa, incluyendo datos de entrada como precio de propiedad, cuota inicial y tasa de interés, así como resultados calculados como cuota fija, VAN y TIR. La relación uno-a-muchos con PaymentScheduleItemModel se establece mediante GORM tags, permitiendo cargar el cronograma completo de pagos en una sola consulta. Los campos calculados se almacenan para optimizar consultas de lectura, siguiendo el patrón CQRS donde comandos y queries pueden tener modelos diferentes.

El PaymentScheduleItemModel, definido en `internal/mortgage/infrastructure/persistence/models/payment_schedule_item_model.go`, representa cada período individual del cronograma de pagos en la tabla "payment_schedule_items". Esta estructura incluye campos para monto de cuota, intereses, amortización y saldo restante, junto con un indicador booleano para identificar períodos de gracia. Los índices compuestos en MortgageID y Period optimizan consultas ordenadas por período, mientras que el índice en UserID facilita filtrado por usuario. El uso de UUID como clave primaria evita conflictos en entornos distribuidos y mejora la seguridad al no exponer secuencias predecibles.

### Repositorio: Operaciones CRUD de Mortgages y Payment Schedule Items

El repositorio de hipotecas, implementado en `internal/mortgage/infrastructure/persistence/repositories/mortgage_repository_impl.go`, maneja operaciones CRUD complejas que involucran tanto la tabla principal de mortgages como la relacionada de payment_schedule_items. El método Save() utiliza transacciones explícitas para garantizar atomicidad, creando primero el registro de mortgage y luego insertando todos los items del cronograma en lote para optimizar rendimiento. El método FindByID() utiliza Preload de GORM para cargar eficientemente la relación uno-a-muchos, reconstruyendo el PaymentSchedule completo en memoria. FindByUserID() implementa paginación con límites y offsets para manejar grandes volúmenes de datos, ordenando resultados por fecha de creación descendente.

El método Update() es particularmente complejo, ya que requiere eliminar todos los items del cronograma anterior y crear nuevos, manteniendo la integridad referencial mediante transacciones. El método Delete() aprovecha las restricciones de clave foránea con CASCADE para eliminar automáticamente los items relacionados, simplificando la lógica de negocio. Los métodos auxiliares toModel() y toDomain() implementan el patrón de mapeo bidireccional, convirtiendo entre entidades de dominio ricas en comportamiento y modelos de persistencia planos optimizados para base de datos. La conversión de tipos incluye validación de value objects, asegurando que datos inválidos no persistan en la base de datos.

## Creación de Tablas: Función autoMigrate()

La función autoMigrate(), definida en `internal/shared/infrastructure/persistence/database.go`, centraliza la creación y migración de todas las tablas de la aplicación. Esta función registra los tres modelos principales - UserModel, MortgageModel y PaymentScheduleItemModel - con GORM, que genera automáticamente las sentencias DDL de SQL basándose en las estructuras de Go y sus tags. La migración automática facilita el desarrollo al sincronizar cambios en el esquema de base de datos con modificaciones en el código, aunque en entornos de producción se recomienda un control más estricto de migraciones mediante scripts versionados.

La configuración de conexión a PostgreSQL incluye parámetros de seguridad como SSL mode y logging apropiado para diferentes entornos. La función NewDatabase() encapsula toda la configuración de conexión, retornando una instancia de GORM lista para uso, siguiendo el patrón de inyección de dependencias para facilitar testing con mocks.

## Conclusión

Esta arquitectura de base de datos demuestra una implementación robusta de DDD con bounded contexts claramente definidos, donde cada contexto mantiene su integridad y responsabilidades específicas. La separación entre IAM y Mortgage permite escalabilidad independiente y reduce riesgos de cambios cruzados, mientras que el uso consistente de repositorios y modelos de persistencia asegura mantenibilidad y testabilidad. Los patrones implementados, como el mapeo entre dominio y persistencia, transacciones para operaciones complejas y optimizaciones de consulta, resultan en un sistema eficiente y confiable para aplicaciones financieras críticas.

## Referencias

Alonso, A. (2018). *Domain-Driven Design con Golang*. Editorial Ra-Ma.

Evans, E. (2003). *Domain-Driven Design: Tackling Complexity in the Heart of Software*. Addison-Wesley.

Vernon, V. (2016). *Implementing Domain-Driven Design*. Addison-Wesley.
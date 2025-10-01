# Calculadora del Método Francés

## Introducción

Esta calculadora implementa el **Método Francés Vencido Ordinario** para calcular cronogramas de pagos de hipotecas. El método francés es un sistema de amortización donde la cuota mensual es fija, compuesta por intereses y amortización del capital. Los intereses se calculan sobre el saldo pendiente al inicio del período.

## Estructura del Proyecto

- **Dominio**: Contiene la lógica de negocio pura (entidades, servicios, repositorios)
- **Aplicación**: Servicios de comandos y consultas
- **Infraestructura**: Persistencia, middleware JWT, etc.
- **Interfaces**: Controladores REST y recursos

La lógica de la calculadora está encapsulada en el servicio `FrenchMethodCalculator` dentro del dominio, siguiendo principios de arquitectura limpia.

## Fórmulas Matemáticas Principales

### 1. Principal Financiado

**Fórmula**: $ P = \text{LoanAmount} - \text{BonoTechoPropio} $

**Explicación**: El principal financiado es el monto del préstamo después de restar cualquier bono o subsidio aplicado. Debe ser mayor que cero.

**Notación Big O**: $ O(1) $ - Operación constante.

### 2. Conversión a Tasa Periódica Mensual

#### Para Tasa Nominal Anual (TNA)
**Fórmula**: $ i = \frac{\text{TNA}}{12} $

**Explicación**: Si la tasa es nominal anual, se divide entre 12 para obtener la tasa mensual efectiva.

#### Para Tasa Efectiva Anual (TEA)
**Fórmula**: $ i = (1 + \text{TEA})^{\frac{1}{12}} - 1 $

**Explicación**: Para tasas efectivas anuales, se calcula la tasa mensual equivalente usando la fórmula de capitalización compuesta.

**Notación Big O**: $ O(1) $ - Operaciones matemáticas constantes.

### 3. Ajuste por Período de Gracia Total

**Fórmula**: $ P_{\text{gracia}} = P \times (1 + i)^{n_{\text{gracia}}} $

**Explicación**: Si hay un período de gracia total, los intereses se capitalizan durante ese período, aumentando el principal que se amortizará posteriormente.

**Notación Big O**: $ O(1) $ - Cálculo exponencial constante.

### 4. Cálculo de la Cuota Fija (Método Francés)

**Fórmula**: $ A = P \times \frac{i \times (1 + i)^n}{(1 + i)^n - 1} $

Donde:
- $ P $: Principal ajustado
- $ i $: Tasa periódica mensual
- $ n $: Número de períodos normales (término total - períodos de gracia)

**Explicación**: Esta es la fórmula clásica del método francés. Calcula la cuota fija que incluye intereses y amortización para pagar el préstamo en n períodos.

**Notación Big O**: $ O(1) $ - Operaciones matemáticas constantes.

### 5. Generación del Cronograma de Pagos

Para cada período k (de 1 a término total):

#### Interés del período
**Fórmula**: $ I_k = \text{saldo}_{k-1} \times i $

#### Amortización
**Fórmula**: $ C_k = A - I_k $ (para períodos normales)

#### Nuevo saldo
**Fórmula**: $ \text{saldo}_k = \text{saldo}_{k-1} - C_k $

#### Tratamiento de períodos de gracia:
- **Gracia Total**: No se paga nada, intereses se capitalizan: $ \text{saldo}_k = \text{saldo}_{k-1} + I_k $
- **Gracia Parcial**: Solo se paga interés: $ \text{saldo}_k = \text{saldo}_{k-1} $

**Notación Big O**: $ O(n) $ donde n es el número total de meses del préstamo. Es lineal porque itera una vez por cada período.

### 6. Valor Actual Neto (VAN/NPV)

**Fórmula**: $ \text{VAN} = -P + \sum_{k=1}^{n} \frac{\text{CF}_k}{(1 + j)^k} $

Donde:
- $ \text{CF}_k $: Flujo de caja en el período k (cuota pagada)
- $ j $: Tasa de descuento mensual
- $ P $: Principal financiado

**Explicación**: El VAN mide la rentabilidad del proyecto descontando todos los flujos de caja a valor presente.

**Notación Big O**: $ O(n) $ - Suma lineal sobre los n períodos.

### 7. Tasa Interna de Retorno (TIR/IRR)

Se calcula usando el **Método de Newton-Raphson** para resolver: $ \text{VAN}(r) = 0 $

**Fórmula iterativa**:
$$
r_{n+1} = r_n - \frac{\text{VAN}(r_n)}{\text{VAN}'(r_n)}
$$

Donde VAN'(r) es la derivada de VAN respecto a r.

**Explicación**: La TIR es la tasa de descuento que hace que el VAN sea cero. Es la tasa de rentabilidad interna del proyecto.

**Notación Big O**: $ O(n \times \text{iteraciones}) $ - Por cada iteración (máximo 1000), se calcula VAN y su derivada, cada una O(n).

### 8. Tasa de Costo Efectivo Anual (TCEA)

**Fórmula**: $ \text{TCEA} = (1 + \text{TIR}_{\text{mensual}})^{12} - 1 $

**Explicación**: Convierte la tasa interna de retorno mensual a su equivalente anual efectivo.

**Notación Big O**: $ O(1) $ - Operación exponencial constante.


### Métodos

- `convertToPeriodicRate()`: Convierte TNA/TEA a tasa mensual
- `calculateFixedInstallment()`: Aplica la fórmula del método francés
- `generatePaymentSchedule()`: Itera por cada período calculando pagos
- `CalculateNPV()`: Computa valor actual neto
- `CalculateIRR()`: Resuelve TIR con Newton-Raphson
- `CalculateTCEA()`: Convierte TIR mensual a anual

## Consideraciones de Rendimiento

- **Complejidad Temporal**: La mayoría de operaciones son O(1), excepto el cronograma O(n) y IRR O(n×iter)
- **Precisión**: Se redondean saldos pequeños (< 0.01) para evitar errores de punto flotante
- **Validaciones**: Se verifican condiciones como principal > 0, términos válidos, etc.

## Conclusión

Esta implementación proporciona una calculadora financiera robusta y precisa del método francés, con todas las fórmulas matemáticas estándar implementadas correctamente. La separación de responsabilidades en capas de dominio/aplicación/infraestructura asegura mantenibilidad y testabilidad.
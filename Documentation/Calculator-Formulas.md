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

**Variables:**

- $ P $: Principal financiado (monto neto del préstamo que se amortizará, excluyendo subsidios)
- $ \text{LoanAmount} $: Monto total del préstamo solicitado por el cliente (monto bruto antes de deducciones)
- $ \text{BonoTechoPropio} $: Bono o subsidio aplicado al "techo propio" (aportación inicial o subsidios gubernamentales que reducen el monto a financiar)

**Notación Big O**: $ O(1) $ - Operación constante.

### 2. Conversión a Tasa Periódica Mensual

#### Para Tasa Nominal Anual (TNA)
**Fórmula**: $ i = \frac{\text{TNA}}{12} $

**Explicación**: Si la tasa es nominal anual, se divide entre 12 para obtener la tasa mensual efectiva.

**Variables:**
- $ i $: Tasa periódica mensual (tasa de interés efectiva mensual)
- $ \text{TNA} $: Tasa Nominal Anual (tasa de interés expresada anualmente, sin considerar capitalización)

#### Para Tasa Efectiva Anual (TEA)
**Fórmula**: $ i = (1 + \text{TEA})^{\frac{1}{12}} - 1 $

**Explicación**: Para tasas efectivas anuales, se calcula la tasa mensual equivalente usando la fórmula de capitalización compuesta.

**Variables:**
- $ i $: Tasa periódica mensual (tasa efectiva mensual equivalente)
- $ \text{TEA} $: Tasa Efectiva Anual (tasa que incluye la capitalización de intereses)

**Notación Big O**: $ O(1) $ - Operaciones matemáticas constantes.

### 3. Ajuste por Período de Gracia Total

**Fórmula**: $ P_{\text{gracia}} = P \times (1 + i)^{n_{\text{gracia}}} $

**Explicación**: Si hay un período de gracia total, los intereses se capitalizan durante ese período, aumentando el principal que se amortizará posteriormente.

**Variables:**
- $ P_{\text{gracia}} $: Principal ajustado después del período de gracia (monto capitalizado)
- $ P $: Principal financiado original
- $ i $: Tasa periódica mensual
- $ n_{\text{gracia}} $: Número de meses de período de gracia total

**Notación Big O**: $ O(1) $ - Cálculo exponencial constante.

### 4. Cálculo de la Cuota Fija (Método Francés)

**Fórmula**: $ A = P \times \frac{i \times (1 + i)^n}{(1 + i)^n - 1} $

**Explicación**: Esta es la fórmula clásica del método francés. Calcula la cuota fija que incluye intereses y amortización para pagar el préstamo en n períodos.

**Variables:**
- $ A $: Cuota fija mensual (monto constante que se paga cada período)
- $ P $: Principal ajustado (monto a amortizar después de ajustes por gracia)
- $ i $: Tasa periódica mensual (tasa de interés efectiva mensual)
- $ n $: Número de períodos normales (término total en meses menos períodos de gracia)

**Notación Big O**: $ O(1) $ - Operaciones matemáticas constantes.

### 5. Generación del Cronograma de Pagos

Para cada período k (de 1 a término total):

#### Interés del período
**Fórmula**: $ I_k = \text{saldo}_{k-1} \times i $

**Variables:**
- $ I_k $: Interés del período k (monto de intereses calculados sobre el saldo pendiente)
- $ \text{saldo}_{k-1} $: Saldo pendiente al inicio del período k (capital restante por amortizar)
- $ i $: Tasa periódica mensual
- $ k $: Número del período actual

#### Amortización
**Fórmula**: $ C_k = A - I_k $ (para períodos normales)

**Variables:**
- $ C_k $: Amortización del período k (parte de la cuota que reduce el capital)
- $ A $: Cuota fija mensual
- $ I_k $: Interés del período k

#### Nuevo saldo
**Fórmula**: $ \text{saldo}_k = \text{saldo}_{k-1} - C_k $

**Variables:**
- $ \text{saldo}_k $: Saldo pendiente al final del período k
- $ \text{saldo}_{k-1} $: Saldo pendiente al inicio del período k
- $ C_k $: Amortización del período k

#### Tratamiento de períodos de gracia:
- **Gracia Total**: No se paga nada, intereses se capitalizan: $ \text{saldo}_k = \text{saldo}_{k-1} + I_k $

  **Variables:**
  - $ \text{saldo}_k $: Saldo capitalizado (incluye intereses no pagados)
  - $ \text{saldo}_{k-1} $: Saldo anterior
  - $ I_k $: Intereses generados (se agregan al capital)

- **Gracia Parcial**: Solo se paga interés: $ \text{saldo}_k = \text{saldo}_{k-1} $

  **Variables:**
  - $ \text{saldo}_k $: Saldo permanece igual (no se amortiza capital)
  - $ \text{saldo}_{k-1} $: Saldo anterior

**Notación Big O**: $ O(n) $ donde n es el número total de meses del préstamo. Es lineal porque itera una vez por cada período.

### 6. Valor Actual Neto (VAN/NPV)

**Fórmula**: $ \text{VAN} = -P + \sum_{k=1}^{n} \frac{\text{CF}_k}{(1 + j)^k} $

**Explicación**: El VAN mide la rentabilidad del proyecto descontando todos los flujos de caja a valor presente.

**Variables:**
- $ \text{VAN} $: Valor Actual Neto (mide si el proyecto es rentable; positivo = rentable)
- $ P $: Principal financiado (inversión inicial, flujo negativo)
- $ \text{CF}_k $: Flujo de caja en el período k (cuotas pagadas por el cliente)
- $ j $: Tasa de descuento mensual (tasa usada para descontar flujos futuros)
- $ n $: Número total de períodos
- $ k $: Índice del período (de 1 a n)

**Notación Big O**: $ O(n) $ - Suma lineal sobre los n períodos.

### 7. Tasa Interna de Retorno (TIR/IRR)

Se calcula usando el **Método de Newton-Raphson** para resolver: $ \text{VAN}(r) = 0 $

**Fórmula iterativa**:
$$
r_{n+1} = r_n - \frac{\text{VAN}(r_n)}{\text{VAN}'(r_n)}
$$

**Explicación**: La TIR es la tasa de descuento que hace que el VAN sea cero. Es la tasa de rentabilidad interna del proyecto.

**Variables:**
- $ r $: Tasa interna de retorno (tasa que hace VAN = 0)
- $ \text{VAN}(r) $: Función del Valor Actual Neto evaluada en r
- $ \text{VAN}'(r) $: Derivada de VAN respecto a r (usada para la aproximación)
- $ n $: Iteración actual del método numérico

**Notación Big O**: $ O(n \times \text{iteraciones}) $ - Por cada iteración (máximo 1000), se calcula VAN y su derivada, cada una O(n).

### 8. Tasa de Costo Efectivo Anual (TCEA)

**Fórmula**: $ \text{TCEA} = (1 + \text{TIR}_{\text{mensual}})^{12} - 1 $

**Explicación**: Convierte la tasa interna de retorno mensual a su equivalente anual efectivo.

**Variables:**
- $ \text{TCEA} $: Tasa de Costo Efectivo Anual (tasa anual efectiva del préstamo)
- $ \text{TIR}_{\text{mensual}} $: Tasa Interna de Retorno mensual (calculada previamente)

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


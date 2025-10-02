# Guía de Parámetros del Endpoint de Cálculo Hipotecario

Esta guía explica cada parámetro del request JSON para el endpoint `/api/v1/mortgage/calculate`, basándose en las fórmulas del **Método Francés** de amortización de préstamos.

---

## Parámetros del Request

### 1. `property_price` (Precio de la Propiedad)

**Tipo:** `number` (mayor a 0)
**Descripción:** Precio total de la vivienda que deseas adquirir.

**Ejemplo:** Si la casa cuesta S/. 150,000, entonces `property_price: 150000`

**Nota:** Este valor se usa principalmente para referencia. La relación típica es:
```
property_price = down_payment + loan_amount
```

---

### 2. `down_payment` (Cuota Inicial o Inicial)

**Tipo:** `number` (mayor o igual a 0)
**Descripción:** Monto que pagas de tu propio bolsillo al momento de comprar la propiedad. Es el "adelanto" o "enganche".

**Fórmula relacionada:**
```
Cuota Inicial (%) = (down_payment / property_price) × 100
```

**Ejemplo:**
- Propiedad de S/. 100,000
- Cuota inicial del 20% = S/. 20,000
- `down_payment: 20000`

**Contexto:** Los bancos generalmente requieren una cuota inicial mínima (10%-30% del valor de la propiedad).

---

### 3. `loan_amount` (Monto del Préstamo Solicitado)

**Tipo:** `number` (mayor a 0)
**Símbolo en fórmulas:** **P₀** (Principal inicial)
**Descripción:** Cantidad de dinero que le solicitas al banco para financiar la compra.

**Fórmula relacionada:**
```
loan_amount = property_price - down_payment
```

**Ejemplo:**
- Propiedad de S/. 100,000
- Cuota inicial de S/. 20,000
- `loan_amount: 80000`

---

### 4. `bono_techo_propio` (Bono del Buen Pagador o Subsidio Estatal)

**Tipo:** `number` (mayor o igual a 0)
**Descripción:** Subsidio otorgado por el gobierno peruano (programa "Techo Propio") para ayudar a familias de bajos recursos a comprar su primera vivienda. Este bono **reduce el monto a financiar**.

**Fórmula del Principal Financiado:**
```
Principal Financiado = loan_amount - bono_techo_propio
```

**Ejemplo:**
- Préstamo solicitado: S/. 80,000
- Bono Techo Propio: S/. 7,500
- Principal que realmente financias: S/. 72,500
- `bono_techo_propio: 7500`

**Nota:** Si no calificas para este subsidio, usa `bono_techo_propio: 0`

---

### 5. `interest_rate` (Tasa de Interés Anual)

**Tipo:** `number` (mayor o igual a 0)
**Símbolo en fórmulas:** **i** (cuando es mensual) o **TNA/TEA** (anual)
**Descripción:** Tasa de interés que el banco cobra por prestarte el dinero. Se expresa en **porcentaje anual** (sin el símbolo %).

**Ejemplo:**
- Si el banco ofrece 7.5% anual, usa `interest_rate: 7.5`
- Si es 12% anual, usa `interest_rate: 12`

**Importante:** El tipo de tasa (nominal o efectiva) se especifica en `rate_type`.

---

### 6. `rate_type` (Tipo de Tasa de Interés)

**Tipo:** `string` (valores permitidos: `"NOMINAL"` o `"EFFECTIVE"`)
**Descripción:** Indica si la tasa de interés proporcionada es **Nominal (TNA)** o **Efectiva (TEA)**.

#### **TNA - Tasa Nominal Anual**
```
Tasa mensual = TNA / 12
```

**Ejemplo:** TNA = 12%
```
Tasa mensual = 12% / 12 = 1% = 0.01
```

#### **TEA - Tasa Efectiva Anual**
```
Tasa mensual = (1 + TEA)^(1/12) - 1
```

**Ejemplo:** TEA = 12%
```
Tasa mensual = (1 + 0.12)^(1/12) - 1 = 0.009489 = 0.9489%
```

**Valores:**
- `"NOMINAL"`: Para TNA
- `"EFFECTIVE"`: Para TEA (más común en el sistema financiero peruano)

---

### 7. `term_months` (Plazo del Préstamo en Meses)

**Tipo:** `integer` (mayor a 0)
**Símbolo en fórmulas:** **n** (número total de períodos)
**Descripción:** Tiempo total en meses que tienes para pagar el préstamo.

**Conversiones comunes:**
- 1 año = 12 meses
- 5 años = 60 meses
- 10 años = 120 meses
- 15 años = 180 meses
- 20 años = 240 meses
- 25 años = 300 meses
- 30 años = 360 meses

**Ejemplo:** Si quieres pagar el préstamo en 20 años, usa `term_months: 240`

---

### 8. `grace_period_months` (Período de Gracia en Meses)

**Tipo:** `integer` (mayor o igual a 0)
**Símbolo en fórmulas:** **n_gracia**
**Descripción:** Número de meses al inicio del préstamo donde **no pagas la cuota normal**. El tipo de gracia se define en `grace_period_type`.

**Ejemplo:** Si el banco te da 6 meses de gracia, usa `grace_period_months: 6`

**Nota:** Si no hay período de gracia, usa `grace_period_months: 0`

---

### 9. `grace_period_type` (Tipo de Período de Gracia)

**Tipo:** `string` (valores permitidos: `"NONE"`, `"PARTIAL"`, `"TOTAL"`)
**Descripción:** Define qué tipo de beneficio recibes durante el período de gracia.

#### **NONE - Sin Gracia**
```
Pagas normalmente desde el primer mes
```
No hay período de gracia. Empiezas a pagar la cuota completa desde el mes 1.

---

#### **PARTIAL - Gracia Parcial**
```
Durante la gracia:
- Cuota = Solo intereses
- Amortización = 0
- Saldo = No disminuye
```

**Ejemplo:** Préstamo de S/. 50,000, TNA 12% (1% mensual), 3 meses de gracia parcial
- Mes 1: Pagas solo S/. 500 (interés)
- Mes 2: Pagas solo S/. 500 (interés)
- Mes 3: Pagas solo S/. 500 (interés)
- Mes 4 en adelante: Pagas la cuota fija normal

El saldo permanece en S/. 50,000 durante los 3 meses de gracia.

---

#### **TOTAL - Gracia Total**
```
Durante la gracia:
- Cuota = 0 (no pagas nada)
- Intereses se capitalizan (se suman al saldo)
- Saldo = Aumenta cada mes
```

**Fórmula del saldo después de gracia total:**
```
Saldo_después = P × (1 + i)^n_gracia
```

**Ejemplo:** Préstamo de S/. 50,000, TEA 10% (TEM 0.7974%), 2 meses de gracia total
- Mes 1: No pagas nada, saldo = S/. 50,398.70
- Mes 2: No pagas nada, saldo = S/. 50,798.99
- Mes 3 en adelante: Pagas cuota fija calculada sobre S/. 50,798.99

**Importante:** La gracia total aumenta el costo total del préstamo porque los intereses se acumulan.

---

### 10. `currency` (Moneda del Préstamo)

**Tipo:** `string` (valores permitidos: `"PEN"` o `"USD"`)
**Descripción:** Moneda en la que se otorga el préstamo.

**Valores:**
- `"PEN"`: Soles peruanos (S/.)
- `"USD"`: Dólares estadounidenses ($)

**Ejemplo:** Para un préstamo en soles, usa `currency: "PEN"`

**Nota:** La moneda no afecta los cálculos matemáticos, solo el símbolo de moneda en los resultados.

---

### 11. `npv_discount_rate` (Tasa de Descuento para el VAN)

**Tipo:** `number` (mayor o igual a 0)
**Símbolo en fórmulas:** **j** (tasa de descuento)
**Descripción:** Tasa de interés anual utilizada para calcular el **Valor Actual Neto (VAN)** del préstamo. Representa el **costo de oportunidad** del dinero.

**Fórmula del VAN:**
```
VAN = -P + Σ[Cuota_k / (1 + j)^k]  para k = 1 hasta n
```

Donde:
- **P** = Principal financiado (desembolso inicial)
- **Cuota_k** = Cuota del período k
- **j** = Tasa de descuento mensual
- **k** = Período

**Ejemplo de uso:**
Si consideras que podrías invertir ese dinero a un 8% anual, usa `npv_discount_rate: 8`

**Interpretación del VAN:**
- VAN > 0: Estás pagando más que el valor presente del dinero
- VAN = 0: El préstamo está en equilibrio
- VAN < 0: Teóricamente favorable (raro en préstamos)

**Nota:** Generalmente se usa una tasa cercana a la tasa de interés del préstamo o la tasa de inflación.

---

## Fórmulas del Método Francés Utilizadas

### 1. Conversión de Tasa Anual a Mensual

**Para TNA (Nominal):**
```
i = TNA / 12
```

**Para TEA (Efectiva):**
```
i = (1 + TEA)^(1/12) - 1
```

---

### 2. Cuota Fija (después del período de gracia)

```
A = P × [i(1 + i)^n] / [(1 + i)^n - 1]
```

Donde:
- **A** = Cuota fija
- **P** = Principal financiado (ajustado si hay gracia total)
- **i** = Tasa de interés mensual
- **n** = Número de períodos normales (term_months - grace_period_months)

---

### 3. Interés del Período k

```
I_k = Saldo_{k-1} × i
```

---

### 4. Amortización del Período k

```
C_k = A - I_k
```

Donde:
- **C_k** = Amortización (pago a capital) del período k
- **A** = Cuota fija
- **I_k** = Interés del período k

---

### 5. Saldo Restante después del Período k

```
Saldo_k = Saldo_{k-1} - C_k
```

---

### 6. Capitalización en Gracia Total

```
Saldo_después_gracia = P × (1 + i)^n_gracia
```

---

### 7. Tasa Interna de Retorno (TIR)

La TIR es la tasa **i** que hace que el VAN = 0:

```
0 = -P + Σ[Cuota_k / (1 + TIR)^k]
```

Se calcula usando el método numérico de Newton-Raphson.

---

### 8. Tasa de Costo Efectivo Anual (TCEA)

```
TCEA = (1 + TIR_mensual)^12 - 1
```

La TCEA representa el **costo real anual** del préstamo, incluyendo todos los pagos.

---

## Ejemplo Completo

```json
{
  "property_price": 150000,
  "down_payment": 22500,
  "loan_amount": 120000,
  "bono_techo_propio": 7500,
  "interest_rate": 7.5,
  "rate_type": "EFFECTIVE",
  "term_months": 240,
  "grace_period_months": 0,
  "grace_period_type": "NONE",
  "currency": "PEN",
  "npv_discount_rate": 8.0
}
```

**Interpretación:**
- Compras una casa de **S/. 150,000**
- Das una cuota inicial de **S/. 22,500** (15%)
- Solicitas un préstamo de **S/. 120,000**
- Recibes un Bono Techo Propio de **S/. 7,500**
- **Principal financiado:** S/. 112,500 (120,000 - 7,500)
- Tasa de interés: **7.5% TEA**
- Plazo: **240 meses (20 años)**
- Sin período de gracia
- Moneda: **Soles (PEN)**
- Tasa de descuento para VAN: **8%**

**Resultados esperados:**
- Tasa mensual: 0.6056%
- Cuota fija: S/. 906.13
- Total de intereses: S/. 104,971.20
- Total a pagar: S/. 217,471.20

---

## Validaciones

El endpoint valida que:

1. `property_price` > 0
2. `down_payment` >= 0
3. `loan_amount` > 0
4. `bono_techo_propio` >= 0
5. `interest_rate` >= 0
6. `rate_type` ∈ {`"NOMINAL"`, `"EFFECTIVE"`}
7. `term_months` > 0
8. `grace_period_months` >= 0
9. `grace_period_type` ∈ {`"NONE"`, `"PARTIAL"`, `"TOTAL"`}
10. `currency` ∈ {`"PEN"`, `"USD"`}
11. `npv_discount_rate` >= 0
12. `loan_amount - bono_techo_propio` > 0 (el principal financiado debe ser positivo)
13. `term_months > grace_period_months` (debe haber períodos normales de pago)

---

## Glosario de Términos Financieros

| Término | Definición |
|---------|------------|
| **Principal** | Monto inicial del préstamo |
| **Interés** | Costo por el uso del dinero prestado |
| **Amortización** | Pago que reduce el saldo del préstamo |
| **Cuota** | Pago mensual total (interés + amortización) |
| **Saldo** | Deuda pendiente en un momento dado |
| **TNA** | Tasa Nominal Anual (simple) |
| **TEA** | Tasa Efectiva Anual (compuesta) |
| **TEM** | Tasa Efectiva Mensual |
| **VAN** | Valor Actual Neto |
| **TIR** | Tasa Interna de Retorno |
| **TCEA** | Tasa de Costo Efectivo Anual (costo real del préstamo) |
| **Período de Gracia** | Tiempo inicial con condiciones especiales de pago |
| **Capitalización** | Suma de intereses no pagados al saldo de la deuda |

---

## Referencias

- Método Francés de Amortización de Préstamos
- Sistema Financiero Peruano - Préstamos Hipotecarios
- Programa Techo Propio - Ministerio de Vivienda, Construcción y Saneamiento del Perú

# Casos de Prueba para el Endpoint `/api/v1/mortgage/calculate`

A continuación se presentan cinco ejemplos de cuerpos **JSON** que pueden usarse para probar la API de cálculo hipotecario con el método francés.
Cada ejemplo incluye una breve **explicación en texto plano**.

---

## Ejemplo 1: Caso típico con Bono Techo Propio

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

**Explicación:**

* Precio de propiedad: S/. 150 000
* Cuota inicial: S/. 22 500 (15 %)
* Bono Techo Propio: S/. 7 500 (subsidio del gobierno)
* Monto a financiar: S/. 120 000
* TEA: 7,5 % (tasa efectiva anual típica en Perú)
* Plazo: 20 años (240 meses)
* Sin período de gracia
* Moneda: **Soles (PEN)**

---

## Ejemplo 2: Con período de gracia **TOTAL**

```json
{
  "property_price": 200000,
  "down_payment": 50000,
  "loan_amount": 140000,
  "bono_techo_propio": 10000,
  "interest_rate": 8.25,
  "rate_type": "EFFECTIVE",
  "term_months": 300,
  "grace_period_months": 12,
  "grace_period_type": "TOTAL",
  "currency": "PEN",
  "npv_discount_rate": 9.0
}
```

**Explicación:**

* Propiedad más cara: S/. 200 000
* Cuota inicial mayor: S/. 50 000 (25 %)
* Bono Techo Propio: S/. 10 000
* 12 meses de gracia **TOTAL** (no paga ni capital ni intereses)
* Plazo: 25 años (300 meses)
* Moneda: **Soles (PEN)**

---

## Ejemplo 3: Con período de gracia **PARCIAL**

```json
{
  "property_price": 180000,
  "down_payment": 36000,
  "loan_amount": 135000,
  "bono_techo_propio": 9000,
  "interest_rate": 7.8,
  "rate_type": "EFFECTIVE",
  "term_months": 180,
  "grace_period_months": 6,
  "grace_period_type": "PARTIAL",
  "currency": "PEN",
  "npv_discount_rate": 8.5
}
```

**Explicación:**

* Precio de propiedad: S/. 180 000
* Cuota inicial: S/. 36 000 (20 %)
* Bono Techo Propio: S/. 9 000
* 6 meses de gracia **PARCIAL** (paga solo intereses, no amortiza capital)
* Plazo: 15 años (180 meses)
* Moneda: **Soles (PEN)**

---

## Ejemplo 4: **Sin Bono Techo Propio**

```json
{
  "property_price": 250000,
  "down_payment": 75000,
  "loan_amount": 175000,
  "bono_techo_propio": 0,
  "interest_rate": 8.9,
  "rate_type": "EFFECTIVE",
  "term_months": 240,
  "grace_period_months": 0,
  "grace_period_type": "NONE",
  "currency": "PEN",
  "npv_discount_rate": 9.5
}
```

**Explicación:**

* Vivienda de mayor valor: S/. 250 000
* Cuota inicial: S/. 75 000 (30 %)
* **Sin subsidio** del gobierno (Bono Techo Propio = 0)
* TEA: 8,9 %
* Plazo: 20 años (240 meses)
* Moneda: **Soles (PEN)**

---

## Ejemplo 5: Con **tasa nominal (TNA)**

```json
{
  "property_price": 120000,
  "down_payment": 24000,
  "loan_amount": 90000,
  "bono_techo_propio": 6000,
  "interest_rate": 7.2,
  "rate_type": "NOMINAL",
  "term_months": 180,
  "grace_period_months": 0,
  "grace_period_type": "NONE",
  "currency": "PEN",
  "npv_discount_rate": 7.5
}
```

**Explicación:**

* Precio de propiedad: S/. 120 000
* Cuota inicial: S/. 24 000 (20 %)
* Bono Techo Propio: S/. 6 000
* **Tasa nominal anual (TNA)**: 7,2 %
* Plazo: 15 años (180 meses)
* Moneda: **Soles (PEN)**

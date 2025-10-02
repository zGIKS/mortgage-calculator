# 🧪 Casos de Prueba - Mortgage API

Endpoint: `POST /api/v1/mortgage/calculate`

---

## ✅ Caso 1 - TEA 12%, 12 meses, sin gracia (PEN)

### Request Body

```json
{
  "property_price": 100000,
  "down_payment": 10000,
  "loan_amount": 90000,
  "bono_techo_propio": 0,
  "interest_rate": 12,
  "rate_type": "EFFECTIVE",
  "term_months": 12,
  "grace_period_months": 0,
  "grace_period_type": "NONE",
  "currency": "PEN",
  "npv_discount_rate": 12
}
```

### Respuesta Esperada

```json
{
  "principal_financed": 90000,
  "periodic_rate": 0.009488792934583046,
  "fixed_installment": 7970.586065049673,
  "total_paid": 95647.0327805961,
  "total_interest_paid": 5647.032780596441,
  "irr": 0.009488792934583046,
  "tcea": 0.12,
  "npv": 0.0,
  "schedule": [
    {
      "period": 1,
      "installment": 7970.586065049673,
      "interest": 853.9913641124741,
      "amortization": 7116.594700937199,
      "balance": 82883.4052990628
    },
    {
      "period": 2,
      "installment": 7970.586065049673,
      "interest": 786.4634705959301,
      "amortization": 7184.122594453744,
      "balance": 75699.28270460905
    },
    "...",
    {
      "period": 11,
      "installment": 7970.586065049673,
      "interest": 149.13645298962638,
      "amortization": 7821.449612060047,
      "balance": 7895.665727877518
    },
    {
      "period": 12,
      "installment": 7970.586065049673,
      "interest": 74.9203371725137,
      "amortization": 7895.66572787716,
      "balance": 0
    }
  ]
}
```

**Tolerancias:** cuotas/intereses ±0.02, saldos ±0.05, tasas ±0.0001

---

## ✅ Caso 2 - TNA 12%, 12 meses, gracia PARCIAL 3 (PEN)

### Request Body

```json
{
  "property_price": 50000,
  "down_payment": 5000,
  "loan_amount": 45000,
  "bono_techo_propio": 0,
  "interest_rate": 12,
  "rate_type": "NOMINAL",
  "term_months": 12,
  "grace_period_months": 3,
  "grace_period_type": "PARTIAL",
  "currency": "PEN",
  "npv_discount_rate": 12
}
```

### Respuesta Esperada

```json
{
  "principal_financed": 45000,
  "periodic_rate": 0.01,
  "fixed_installment": 5253.316328235639,
  "total_paid": 48629.84695412075,
  "total_interest_paid": 3629.8469541208497,
  "irr": 0.01,
  "tcea": 0.12682503013196977,
  "npv": 177.15835489947312,
  "schedule": [
    {
      "period": 1,
      "installment": 450,
      "interest": 450,
      "amortization": 0,
      "balance": 45000,
      "grace_type": "PARTIAL"
    },
    {
      "period": 2,
      "installment": 450,
      "interest": 450,
      "amortization": 0,
      "balance": 45000,
      "grace_type": "PARTIAL"
    },
    {
      "period": 3,
      "installment": 450,
      "interest": 450,
      "amortization": 0,
      "balance": 45000,
      "grace_type": "PARTIAL"
    },
    {
      "period": 4,
      "installment": 5253.316328235639,
      "interest": 450,
      "amortization": 4803.316328235639,
      "balance": 40196.68367176436
    },
    "...",
    {
      "period": 12,
      "installment": 5253.316328235639,
      "interest": 52.01303295282907,
      "amortization": 5201.30329528281,
      "balance": 0
    }
  ]
}
```

**Tolerancias:** cuotas/intereses ±0.02, saldos ±0.05, tasas ±0.0001

---

## ⚠️ Caso 3 - TEA 10%, 6 meses, gracia TOTAL 2, con bono (PEN) - CORREGIDO

**Nota:** El valor de `total_interest_paid` fue corregido de **508.37** a **908.67** porque debe incluir los intereses capitalizados durante la gracia TOTAL.

### Request Body

```json
{
  "property_price": 30000,
  "down_payment": 3000,
  "loan_amount": 27000,
  "bono_techo_propio": 2000,
  "interest_rate": 10,
  "rate_type": "EFFECTIVE",
  "term_months": 6,
  "grace_period_months": 2,
  "grace_period_type": "TOTAL",
  "currency": "PEN",
  "npv_discount_rate": 10
}
```

### Respuesta Esperada

```json
{
  "principal_financed": 25000,
  "periodic_rate": 0.007974140428903007,
  "capitalized_balance_after_grace": 25400.29670924513,
  "fixed_installment": 6477.167849710664,
  "total_paid": 25908.67139884266,
  "total_interest_paid": 908.67139884266,
  "irr": 0.007974140428903007,
  "tcea": 0.10,
  "npv": 0.0,
  "schedule": [
    {
      "period": 1,
      "installment": 0,
      "interest": 199.35351072257518,
      "amortization": 0,
      "balance": 25199.353510722575,
      "grace_type": "TOTAL",
      "interest_capitalized": true
    },
    {
      "period": 2,
      "installment": 0,
      "interest": 200.94319852255502,
      "amortization": 0,
      "balance": 25400.29670924513,
      "grace_type": "TOTAL",
      "interest_capitalized": true
    },
    {
      "period": 3,
      "installment": 6477.167849710664,
      "interest": 202.5402967092451,
      "amortization": 6274.627552,
      "balance": 19125.669156
    },
    "...",
    {
      "period": 6,
      "installment": 6477.167849710664,
      "interest": 51.87,
      "amortization": 6425.297,
      "balance": 0
    }
  ]
}
```

**Tolerancias:** cuotas/intereses ±0.02, saldos ±0.10, tasas ±0.0001

**✅ Corrección aplicada:** `total_interest_paid = 908.67` (incluye intereses capitalizados de ~400.30 + intereses pagados de ~508.37)

---

## ✅ Caso 4 - TEA 7.5%, 240 meses, con bono (PEN)

### Request Body

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

### Respuesta Esperada

```json
{
  "principal_financed": 112500,
  "periodic_rate": 0.006044919024291717,
  "fixed_installment": 889.4390329991267,
  "total_paid": 213465.3679197909,
  "total_interest_paid": 100965.36791978808,
  "irr": 0.006044919024291717,
  "tcea": 0.075,
  "npv": -5500.0,
  "schedule": [
    {
      "period": 1,
      "installment": 889.4390329991267,
      "interest": 680.0533902328182,
      "amortization": 209.38564276630848,
      "balance": 112290.61435723369
    },
    {
      "period": 2,
      "installment": 889.4390329991267,
      "interest": 678.7876709774465,
      "amortization": 210.65136202168014,
      "balance": 112079.962995212
    },
    "...",
    {
      "period": 239,
      "installment": 889.4390329991267,
      "interest": 10.65645073425536,
      "amortization": 878.7825822648713,
      "balance": 884.0947518122341
    },
    {
      "period": 240,
      "installment": 889.4390329991267,
      "interest": 5.344281184506238,
      "amortization": 884.0947518146204,
      "balance": 0
    }
  ]
}
```

**Tolerancias:** cuotas/intereses ±0.02, saldos ±0.10, tasas ±0.0001

**Nota:** NPV negativo porque `npv_discount_rate` (8%) > `interest_rate` (7.5%)

---

## ✅ Caso 5 - USD, TEA 8%, 12 meses, sin gracia

### Request Body

```json
{
  "property_price": 80000,
  "down_payment": 16000,
  "loan_amount": 64000,
  "bono_techo_propio": 0,
  "interest_rate": 8,
  "rate_type": "EFFECTIVE",
  "term_months": 12,
  "grace_period_months": 0,
  "grace_period_type": "NONE",
  "currency": "USD",
  "npv_discount_rate": 8
}
```

### Respuesta Esperada

```json
{
  "principal_financed": 64000,
  "periodic_rate": 0.006434030109,
  "fixed_installment": 5559.002045,
  "total_paid": 66708.02454,
  "total_interest_paid": 2708.02454,
  "irr": 0.006434030109,
  "tcea": 0.08,
  "npv": 0.0
}
```

**Tolerancias:** cuotas/intereses ±0.02, saldos ±0.05, tasas ±0.0001

---

## ⚠️ Caso 6 - USD, TNA 9.6%, 18 meses, gracia PARCIAL 6 - CORREGIDO

**Nota:** El valor de `tcea` fue corregido de **0.099556** a **0.1003386937** según la fórmula correcta: (1 + 0.008)^12 - 1

### Request Body

```json
{
  "property_price": 50000,
  "down_payment": 10000,
  "loan_amount": 40000,
  "bono_techo_propio": 0,
  "interest_rate": 9.6,
  "rate_type": "NOMINAL",
  "term_months": 18,
  "grace_period_months": 6,
  "grace_period_type": "PARTIAL",
  "currency": "USD",
  "npv_discount_rate": 10
}
```

### Respuesta Esperada

```json
{
  "principal_financed": 40000,
  "periodic_rate": 0.008,
  "fixed_installment": 3509.198441,
  "total_paid": 44030.3806,
  "total_interest_paid": 4030.3806,
  "irr": 0.008,
  "tcea": 0.1003386937,
  "npv": -150.0,
  "schedule": [
    {
      "period": 1,
      "installment": 320,
      "interest": 320,
      "amortization": 0,
      "balance": 40000,
      "grace_type": "PARTIAL"
    },
    {
      "period": 2,
      "installment": 320,
      "interest": 320,
      "amortization": 0,
      "balance": 40000,
      "grace_type": "PARTIAL"
    },
    "...",
    {
      "period": 6,
      "installment": 320,
      "interest": 320,
      "amortization": 0,
      "balance": 40000,
      "grace_type": "PARTIAL"
    },
    {
      "period": 7,
      "installment": 3509.198441,
      "interest": 320,
      "amortization": 3189.198441,
      "balance": 36810.801559
    },
    "...",
    {
      "period": 18,
      "installment": 3509.198441,
      "interest": 27.27,
      "amortization": 3481.93,
      "balance": 0
    }
  ]
}
```

**Tolerancias:** cuotas/intereses ±0.02, saldos ±0.05, tasas ±0.0001

**✅ Corrección aplicada:** `tcea = 0.1003386937` (fórmula: (1.008)^12 - 1 = 0.10034)

---

## 📋 Reglas de Validación

### Identidad matemática
```
installment = interest + amortization (en cada período)
```

### Cierre del préstamo
```
balance[último_período] ≈ 0 (±0.05)
suma_amortizaciones ≈ principal_financed
```

### NPV (Valor Presente Neto)
```
Si npv_discount_rate == interest_rate (TEA) → NPV ≈ 0
Si npv_discount_rate > interest_rate → NPV < 0
Si npv_discount_rate < interest_rate → NPV > 0
```

### TCEA (Tasa de Costo Efectiva Anual)
```
tcea = (1 + irr_mensual)^12 - 1
```

---

## ⚙️ Tolerancias Recomendadas

| Métrica | Tolerancia |
|---------|-----------|
| Cuotas, intereses, amortizaciones | ±0.02 |
| Saldos (plazos cortos) | ±0.05 |
| Saldos (plazos largos >120 meses) | ±0.10 |
| Tasas (IRR, TCEA, periodic_rate) | ±0.0001 |
| Saldo final | ±0.05 |
| NPV | ±1.0 (±10.0 para plazos largos) |

---

## 🔧 Correcciones Aplicadas

1. **Caso 3:** `total_interest_paid` corregido de **508.37** a **908.67**
   - Razón: Debe incluir intereses capitalizados en gracia TOTAL

2. **Caso 6:** `tcea` corregido de **0.099556** a **0.1003386937**
   - Razón: Aplicación correcta de la fórmula (1 + i_mensual)^12 - 1
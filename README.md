# API

## change log
- v3_payments增加evname字段

```sql
ALTER TABLE v3_payments ADD COLUMN evName TEXT ;
UPDATE v3_payments set evName="" where evName is null;
```

- coin icon url
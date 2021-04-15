# API

## Install
```shell
git clone https://github.com/wolot/api.git
cd api && make
cd bin && ./api
```

## change log
- v3_payments增加evname字段

```sql
ALTER TABLE v3_payments ADD COLUMN evName TEXT ;
UPDATE v3_payments set evName="" where evName is null;
```

- coin icon url
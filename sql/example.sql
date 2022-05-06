---
---
--- After running the go program check content of transaction table
---
select t.id, af.number from_account, tof.number to_account, d.name categorie 
from transaction t, account af, account tof, target d 
where  t.target = d.id and t.from_account=af.id and t.to_account=tof.id;

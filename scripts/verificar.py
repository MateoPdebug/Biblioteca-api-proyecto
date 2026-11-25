"""Pruebas HTTP reales. Ejecutar con Python 3 y Docker Compose activo."""
import json, subprocess, urllib.request, urllib.error, uuid
from concurrent.futures import ThreadPoolExecutor
BASE='http://localhost:8080'
run=uuid.uuid4().hex[:10]
def req(method,path,status,body=None,token=None):
    headers={'Content-Type':'application/json'}
    if token: headers['Authorization']='Bearer '+token
    q=urllib.request.Request(BASE+path,data=json.dumps(body).encode() if body is not None else None,headers=headers,method=method)
    try:
        with urllib.request.urlopen(q,timeout=15) as r: code=r.status; data=json.load(r)
    except urllib.error.HTTPError as e: code=e.code; data=json.load(e)
    print(f'{method} {path}: HTTP {code} (esperado {status})',flush=True)
    assert code==status,(path,code,data)
    if 'token' not in data: print(json.dumps(data,ensure_ascii=False),flush=True)
    return data

def register(label):
    email=f'{label}-{run}@example.com'; password='Prueba123!'
    req('POST','/api/auth/register',201,dict(name=label,email=email,password=password))
    return email,password,req('POST','/api/auth/login',200,dict(email=email,password=password))['token']

assert req('GET','/health',200)=={'status':'ok'}
assert req('GET','/ping',200)=={'message':'pong'}
_,_,student=register('student'); _,_,other=register('other'); email,pwd,_=register('admin')
subprocess.run(['docker','compose','exec','-T','mysql','mysql','-uroot','-proot_pass','biblioteca','-e',f"UPDATE users SET role='admin' WHERE email='{email}'"],check=True)
admin=req('POST','/api/auth/login',200,dict(email=email,password=pwd))['token']
books=req('GET','/api/books',200,token=student)
book=next(b for b in books if b['available_copies']>0)
loan=req('POST','/api/loans',201,{'book_id':book['id']},student)
req('GET','/api/loans/me',200,token=student)
req('PUT',f"/api/loans/{loan['id']}/return",200,token=student)
body=dict(serial='TEST-'+run,brand='Dell',model='Latitude')
req('GET','/api/laptops',401)
req('POST','/api/laptops',403,body,student)
laptop=req('POST','/api/laptops',201,body,admin)
req('POST','/api/laptops',409,body,admin)
id=laptop['id']; payload={'laptop_id':id}
assert any(x['id']==id and x['available'] for x in req('GET','/api/laptops',200,token=student))
l=req('POST','/api/laptop-loans',201,payload,student)
req('POST','/api/laptop-loans',409,payload,other)
assert any(x['id']==l['id'] for x in req('GET','/api/laptop-loans/me',200,token=student))
assert not any(x['id']==l['id'] for x in req('GET','/api/laptop-loans/me',200,token=other))
path=f"/api/laptop-loans/{l['id']}/return"
req('PUT',path,403,token=other)
req('PUT',path,200,token=student)
req('PUT',path,409,token=student)
assert any(x['id']==id and x['available'] for x in req('GET','/api/laptops',200,token=student))
l=req('POST','/api/laptop-loans',201,payload,other)
req('PUT',f"/api/laptop-loans/{l['id']}/return",200,token=admin)
# Dos usuarios intentan adquirir el mismo equipo simultáneamente.
def race(token):
    q=urllib.request.Request(BASE+'/api/laptop-loans',data=json.dumps(payload).encode(),headers={'Content-Type':'application/json','Authorization':'Bearer '+token},method='POST')
    try:
        with urllib.request.urlopen(q,timeout=15) as r: return r.status,json.load(r),token
    except urllib.error.HTTPError as e: return e.code,json.load(e),token
with ThreadPoolExecutor(max_workers=2) as pool:
    results=list(pool.map(race,[student,other]))
print('Préstamos simultáneos:', [r[0] for r in results],flush=True)
assert sorted(r[0] for r in results)==[201,409],results
winner=next(r for r in results if r[0]==201)
req('PUT',f"/api/laptop-loans/{winner[1]['id']}/return",200,token=winner[2])
print('TODAS LAS PRUEBAS HTTP SUPERADAS')

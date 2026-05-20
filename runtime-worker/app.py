from fastapi import FastAPI
from fastapi.responses import StreamingResponse
app=FastAPI()
@app.get('/health')
def h(): return {'ok':True}
@app.get('/metrics')
def m(): return 'runtime_up 1'
@app.get('/v1/models')
def models(): return {'data':[{'id':'Qwen/Qwen2.5-7B-Instruct'}]}
@app.post('/v1/chat/completions')
def c(): return {'id':'rt-1','choices':[{'message':{'role':'assistant','content':'hello'}}]}

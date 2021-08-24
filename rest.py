from aiohttp import web
import schedule
import time
import threading
running=True



async def handle(request):

    job_thread= threading.Thread(target = init)
    job_thread.start() 
    return web.Response(text="hola")


async def stop(request):
    print("llego aqui")
    global running
    running= False
    return web.Response(text="stoped")



def job():
    print("hola")

def init():
    schedule.every(0.05).seconds.do(job)
    global running
    while running:
        print(running)
        schedule.run_pending()
        time.sleep(0.05)
    running=True
    return 

# 
app= web.Application()
app.add_routes([
    web.get('/',handle),
    web.get('/stop',stop)
])

if __name__=='__main__':
    web.run_app(app,port=8585)


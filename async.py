import asyncio
import motor.motor_asyncio



def get_db():
    client=motor.motor_asyncio.AsyncIOMotorClient('localhost',27017)

    try:
        
        return client.test2
    except Exception:
        print("Unable to connect to the server")

async def insertData(db):

    collection=db.links
    document={'id':'d1'}
    result= await collection.insert_one(document)
    print('result %s',repr(result.inserted_id))


db=get_db()
loop= asyncio.get_event_loop()
loop.run_until_complete(insertData(db))
print("hola mundo")
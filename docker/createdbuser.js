
db.createUser({
	user: 'neurone', pwd: 'neur0n3',
	roles: [{ role: "readWrite", db: "neurone" }]
});


db.getCollection('sequences').insert([
	{ actual_value: 0 }
]);


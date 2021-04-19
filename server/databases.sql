CREATE TABLE IF NOT EXISTS pokemons (
    id INTEGER PRIMARY KEY, 
    name TEXT NOT NULL,
    life INTEGER NOT NULL,
    type TEXT NOT NULL, 
    level INTEGER NOT NULL 
);


INSERT INTO pokemons (name,life,type,level) 
    VALUES 
        ('pikachu',200,'electric',32),
        ('charizard',300,'fire',64),
        ('seal',250,'water',64),
        ('electabuzz',260,'electric',12),
        ('flareon',250,'fire',8),
        ('mewtwo',500,'psychic',100),
        ('mew',500,'psychic',100);
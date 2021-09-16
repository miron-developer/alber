-- future feature
CREATE TABLE IF not EXISTS Countries (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE
);

-- country codes for mobile
CREATE TABLE IF not EXISTS CountryCodes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL,
    countryID INTEGER NOT NULL,
    FOREIGN KEY (countryID) REFERENCES Countries(id) ON DELETE CASCADE
);

-- when you need country
-- create migrations with change
CREATE TABLE IF not EXISTS Cities (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    countryID INTEGER NOT NULL,
    FOREIGN KEY (countryID) REFERENCES Countries(id) ON DELETE CASCADE
);

-- travelTypes
-- 1 = car
-- 2 = train
-- 3 = plane
-- 4 = ship
CREATE TABLE IF not EXISTS TravelTypes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE
);

-- topTypes
-- 1 - 1 day
-- 2 - 3 days
-- 3 - 5 days
-- 4 - 1 week
-- 5 - 10 days
-- 6 - 2 weeks
-- 7 - 1 month (equal 30 days)
-- 8 - 2 month
-- going futher is not meaningful
CREATE TABLE IF not EXISTS TopTypes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE,
    color TEXT UNIQUE, -- you must write correct colour)
    duration INTEGER NOT NULL,
    cost INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS Users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	nickname TEXT UNIQUE,
    phoneNumber TEXT UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE IF not EXISTS Sessions (
	id TEXT PRIMARY KEY,
    expireDatetime INTEGER NOT NULL,
    userID INTEGER UNIQUE,
	FOREIGN KEY (userID) REFERENCES Users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Parsels (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    title TEXT NOT NULL,
    weight INTEGER NOT NULL,
    price INTEGER NOT NULL,
    contactNumber TEXT NOT NULL,
    creationDatetime INTEGER NOT NULL,
    expireDatetime INTEGER NOT NULL,
    expireOnTopDatetime INTEGER,
    isHaveWhatsUp INTEGER NOT NULL,
    userID INTEGER NOT NULL,
    fromID INTEGER NOT NULL,
    toID INTEGER NOT NULL,
    topTypeID INTEGER,
    FOREIGN KEY (userID) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (topTypeID) REFERENCES TopTypes(id) ON DELETE CASCADE,
    FOREIGN KEY (fromID) REFERENCES Cities(id) ON DELETE CASCADE,
    FOREIGN KEY (toID) REFERENCES Cities(id) ON DELETE CASCADE,
    CHECK(
        LENGTH(title) <= 100 AND
        isHaveWhatsUp IN (0, 1) AND (
            (creationDatetime < expireDatetime) OR
            (creationDatetime < expireOnTopDatetime AND expireOnTopDatetime <= expireDatetime AND expireOnTopDatetime IS NOT NULL)
        ) AND
        fromID != toID
    )
);

CREATE TABLE IF NOT EXISTS Travelers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    contactNumber TEXT NOT NULL,
    weight INTEGER NOT NULL,
    creationDatetime INTEGER NOT NULL,
    departureDatetime INTEGER NOT NULL,
    arrivalDatetime INTEGER NOT NULL,
    expireOnTopDatetime INTEGER,
    isHaveWhatsUp INTEGER NOT NULL,
    userID INTEGER NOT NULL,
    travelTypeID INTEGER NOT NULL,
    fromID INTEGER NOT NULL,
    toID INTEGER NOT NULL,
    topTypeID INTEGER,
    FOREIGN KEY (userID) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (travelTypeID) REFERENCES TravelTypes(id) ON DELETE CASCADE,
    FOREIGN KEY (topTypeID) REFERENCES TopTypes(id) ON DELETE CASCADE,
    FOREIGN KEY (fromID) REFERENCES Cities(id) ON DELETE CASCADE,
    FOREIGN KEY (toID) REFERENCES Cities(id) ON DELETE CASCADE,
    CHECK(
        isHaveWhatsUp IN (0, 1) AND (
            (creationDatetime < departureDatetime AND departureDatetime < arrivalDatetime) OR
            (creationDatetime < expireOnTopDatetime AND expireOnTopDatetime <= departureDatetime AND departureDatetime < arrivalDatetime AND expireOnTopDatetime IS NOT NULL) 
        ) AND
        fromID != toID
    )
);

-- clipped images
CREATE TABLE IF NOT EXISTS Images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    src TEXT NOT NULL,
    name TEXT NOT NULL,
    userID INTEGER,
    parselID INTEGER,
    FOREIGN KEY (userID) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (parselID) REFERENCES Parsels(id) ON DELETE CASCADE
);
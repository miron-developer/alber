-- create a new table 
CREATE TABLE IF NOT EXISTS Travelers2 (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL,
    contactNumber TEXT NOT NULL,
    creationDatetime INTEGER NOT NULL,
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
        isHaveWhatsUp IN (0, 1) AND
        fromID != toID AND
        LENGTH(description) < 1000
    )
);
-- copy data from old table to the new one
INSERT INTO Travelers2 SELECT id, weight, contactNumber, creationDatetime, expireOnTopDatetime, isHaveWhatsUp, userID, travelTypeID, fromID, toID, topTypeID FROM Travelers;

-- drop the old table
DROP TABLE Travelers;

-- rename new table to the old one
ALTER TABLE Travelers2 RENAME TO Travelers;
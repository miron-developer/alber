-- create a new table 
CREATE TABLE IF NOT EXISTS Parsels2 (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    description TEXT NOT NULL,
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
        LENGTH(description) <= 1000 AND 
        isHaveWhatsUp IN (0, 1) AND (
            (creationDatetime < expireDatetime) OR
            (creationDatetime < expireOnTopDatetime AND expireOnTopDatetime <= expireDatetime AND expireOnTopDatetime IS NOT NULL)
        ) AND
        fromID != toID
    )
);
-- copy data from old table to the new one
INSERT INTO Parsels2 SELECT * FROM Parsels;

-- drop the old table
DROP TABLE Parsels;

-- rename new table to the old one
ALTER TABLE Parsels2 RENAME TO Parsels;
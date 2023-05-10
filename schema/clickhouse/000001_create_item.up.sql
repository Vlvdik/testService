CREATE TABLE IF NOT EXISTS items
(
    `Id` Int32,
    `CampaignId` Int32,
    `Name` String,
    `Description` String,
    `Priority` Int32,
    `Removed` Bool,
    `EventTime` DateTime
)
    ENGINE=Log;

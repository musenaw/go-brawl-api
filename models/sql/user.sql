CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  tag TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  nameColor TEXT,
  trophies INTEGER DEFAULT 0,
  highestTrophies INTEGER DEFAULT 0,
  highestPowerPlayPoints INTEGER DEFAULT 0,
  expLevel INTEGER DEFAULT 0,
  expPoints INTEGER DEFAULT 0,
  isQualifiedFromChampionshipChallenge INTEGER DEFAULT 0,
  teamVictories INTEGER DEFAULT 0,
  soloVictories INTEGER DEFAULT 0,
  duoVictories INTEGER DEFAULT 0,
  bestRoboRumbleTime INTEGER DEFAULT 0,
  bestTimeAsBigBrawler INTEGER DEFAULT 0
);

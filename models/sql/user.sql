CREATE TABLE users (
  ID SERIAL PRIMARY KEY,
  tag TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  name_color TEXT,
  trophies INTEGER DEFAULT 0,
  highest_trophies INTEGER DEFAULT 0,
  highest_power_playPoints INTEGER DEFAULT 0,
  exp_level INTEGER DEFAULT 0,
  exp_points INTEGER DEFAULT 0,
  is_qualified_from_championship_challenge BOOLEAN,
  team_victories INTEGER DEFAULT 0,
  solo_victories INTEGER DEFAULT 0,
  duoVictories INTEGER DEFAULT 0,
  best_robo_rumble_time INTEGER DEFAULT 0,
  best_time_as_big_brawler INTEGER DEFAULT 0
);

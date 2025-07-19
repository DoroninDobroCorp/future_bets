-- +migrate Up
CREATE TABLE IF NOT EXISTS analyzer.leagues(
	id BIGSERIAL PRIMARY KEY,
  	bookmaker_name VARCHAR(30) NOT NULL,
  	sport_name VARCHAR(30) NOT NULL,
  	league_name VARCHAR(100) NOT NULL,
	created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  	UNIQUE(bookmaker_name, sport_name, league_name)
);

CREATE TABLE IF NOT EXISTS analyzer.teams(
	id BIGSERIAL PRIMARY KEY,
  	league_id BIGINT NOT NULL,
  	team_name VARCHAR(100) NOT NULL,
	created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  	UNIQUE(league_id, team_name),
  	FOREIGN KEY (league_id) REFERENCES analyzer.leagues(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS analyzer.leagues_merge(
	id BIGSERIAL PRIMARY KEY,
  	league1_id BIGINT,
  	league2_id BIGINT,
	created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  	FOREIGN KEY (league1_id) REFERENCES analyzer.leagues(id) ON DELETE CASCADE,
  	FOREIGN KEY (league2_id) REFERENCES analyzer.leagues(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS analyzer.teams_merge(
	uuid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  	team1_id BIGINT,
  	team2_id BIGINT,
	created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    FOREIGN KEY (team1_id) REFERENCES analyzer.teams(id) ON DELETE CASCADE,
  	FOREIGN KEY (team2_id) REFERENCES analyzer.teams(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE analyzer.teams_merge;
DROP TABLE analyzer.leagues_merge;
DROP TABLE analyzer.teams;
DROP TABLE analyzer.leagues;
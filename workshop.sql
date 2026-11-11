CREATE TABLE published_files (
	id SERIAL,
	publishedfileid BIGINT PRIMARY KEY,
	title TEXT NOT NULL,
	time_created BIGINT NOT NULL,
	time_updated BIGINT NOT NULL,
	creator BIGINT NOT NULL,
	tags text[] NOT NULL,
	file_size BIGINT NOT NULL,
	file_url TEXT NOT NULL,
	preview_url TEXT NOT NULL,
	subscriptions BIGINT NOT NULL,
	creator_appid BIGINT NOT NULL,
	consumer_appid BIGINT NOT NULL,
	detail JSONB NOT NULL
);

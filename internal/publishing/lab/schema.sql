CREATE SCHEMA f09_receiver;
CREATE TABLE f09_receiver.inbox (
 dataset text NOT NULL, owner text NOT NULL, event_id text NOT NULL,
 order_id text NOT NULL, seq bigint NOT NULL CHECK (seq>0), kind text NOT NULL, body bytea NOT NULL,
 PRIMARY KEY(dataset,owner,event_id), UNIQUE(dataset,owner,order_id,seq)
);
CREATE TABLE f09_receiver.job (
 dataset text NOT NULL, owner text NOT NULL, event_id text NOT NULL, done boolean NOT NULL DEFAULT false,
 PRIMARY KEY(dataset,owner,event_id), FOREIGN KEY(dataset,owner,event_id) REFERENCES f09_receiver.inbox
);
CREATE TABLE f09_receiver.projection (
 dataset text NOT NULL, owner text NOT NULL, order_id text NOT NULL,
 cursor bigint NOT NULL DEFAULT 0, state text NOT NULL DEFAULT '', effects bigint NOT NULL DEFAULT 0,
 PRIMARY KEY(dataset,owner,order_id)
);
CREATE TABLE f09_receiver.diagnostic (
 topic text NOT NULL, part integer NOT NULL, pos bigint NOT NULL, reason text NOT NULL,
 PRIMARY KEY(topic,part,pos)
);

CREATE INDEX idx_music_artis ON stream.music_artis(music_id);
CREATE UNIQUE INDEX idx_artis ON stream.artis(artis_id, slug);
create table if not exists book
(
    id    serial primary key,
    title varchar NOT NULL unique
);

create table if not exists answer
(
    id             serial primary key,
    book_id        int references book (id),
    section_type   varchar check ( section_type in ('READING', 'LISTENING')),
    section_answer TEXT[] NOT NULL,
    UNIQUE (book_id, section_type)
);

create table if not exists exam
(
    id                  uuid primary key,
    user_id             int                      NOT NULL,
    book_id             int references book (id) NOT NULL,
    over_all_band_score int                      NOT NULL DEFAULT 0
);


create table if not exists speaking_detail
(
    id               uuid primary key,
    exam_id          uuid references exam (id) NOT NULL,
    part_number      int check ( part_number in (1, 2, 3)),
    fluency_score    int check ( fluency_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    grammar_score    int check ( grammar_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    vocabulary_score int check ( vocabulary_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    coherence_score  int check ( coherence_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    topic_dev_score  int check ( coherence_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    relevance_score  int check ( coherence_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    sentence_count   int                       NOT NULL DEFAULT 0,
    word_count       int                       NOT NULL DEFAULT 0,
    transcription    jsonb,
    voice_url        varchar,
    speaking_score   int                       NOT NULL DEFAULT 0
);


create table if not exists writing_detail
(
    id                     uuid primary key,
    exam_id                uuid references exam (id)          NOT NULL,
    task_number            int check ( task_number in (1, 2)) NOT NULL,
    response               jsonb                              NOT NULL,
    feedback               TEXT,
    coherence_score        int check (coherence_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    grammar_score          int check ( grammar_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    lexical_resource_score int check ( lexical_resource_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    task_achievement_score int check ( task_achievement_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9) ),
    task_score             int                                NOT NULL DEFAULT 0
);

create table if not exists listening_detail
(
    id          uuid primary key,
    exam_id     uuid references exam (id) NOT NULL,
    band_score  int                       NOT NULL DEFAULT 0 check ( band_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    user_answer jsonb                     NOT NULL
);

create table if not exists reading_detail
(
    id          uuid primary key,
    exam_id     uuid references exam (id) NOT NULL,
    band_score  int                       NOT NULL DEFAULT 0 check ( band_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    user_answer jsonb                     NOT NULL
);
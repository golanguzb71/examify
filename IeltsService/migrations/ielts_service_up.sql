create table if not exists book
(
    id         serial primary key,
    title      varchar NOT NULL unique,
    created_at timestamp DEFAULT now()
);

create table if not exists answer
(
    id             serial primary key,
    book_id        int references book (id),
    section_type   varchar check ( section_type in ('READING', 'LISTENING', 'WRITING')),
    section_answer TEXT[] NOT NULL,
    created_at     timestamp DEFAULT now(),
    UNIQUE (book_id, section_type)
);

create table if not exists exam
(
    id                  uuid primary key,
    user_id             int                      NOT NULL,
    book_id             int references book (id) NOT NULL,
    over_all_band_score int                      NOT NULL DEFAULT 0,
    created_at          timestamp                         DEFAULT now()
);


create table if not exists speaking_detail
(
    id               uuid primary key,
    exam_id          uuid references exam (id) NOT NULL,
    part_number      int check ( part_number in (1, 2, 3)),
    fluency_score    float check ( fluency_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    grammar_score    float check ( grammar_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    vocabulary_score float check ( vocabulary_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    coherence_score  float check ( coherence_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    topic_dev_score  float check ( topic_dev_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    relevance_score  float check ( relevance_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    sentence_count   int                       NOT NULL DEFAULT 0,
    word_count       int                       NOT NULL DEFAULT 0,
    transcription    jsonb,
    voice_url        varchar,
    speaking_score   float                       NOT NULL DEFAULT 0,
    created_at       timestamp                          DEFAULT now()
);


create table if not exists writing_detail
(
    id                     uuid primary key,
    exam_id                uuid references exam (id)          NOT NULL,
    task_number            int check ( task_number in (1, 2)) NOT NULL,
    response               jsonb                              NOT NULL,
    feedback               TEXT,
    coherence_score        float check (coherence_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    grammar_score          float check ( grammar_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    lexical_resource_score float check ( lexical_resource_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    task_achievement_score float check ( task_achievement_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9) ),
    task_score             float                                NOT NULL DEFAULT 0,
    created_at             timestamp                                   DEFAULT now()
);

create table if not exists listening_detail
(
    id          uuid primary key,
    exam_id     uuid references exam (id) NOT NULL UNIQUE,
    band_score  float                       NOT NULL DEFAULT 0 check ( band_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    user_answer jsonb                     NOT NULL,
    created_at  timestamp                          DEFAULT now()
);

create table if not exists reading_detail
(
    id          uuid primary key,
    exam_id     uuid references exam (id) NOT NULL UNIQUE,
    band_score  float                       NOT NULL DEFAULT 0 check ( band_score in (0, 1, 2, 3, 4, 5, 6, 7, 8, 9)),
    user_answer jsonb                     NOT NULL,
    created_at  timestamp                          DEFAULT now()
);
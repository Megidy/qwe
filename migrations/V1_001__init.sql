CREATE TABLE cats(
    id UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL,

    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 

    on_mission BOOL NOT NULL DEFAULT false,

    breed VARCHAR(100) NOT NULL,
    salary NUMERIC(10, 2) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE missions(
    id UUID PRIMARY KEY,
    cat_id UUID,

    is_completed BOOL NOT NULL DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_cat FOREIGN KEY (cat_id) REFERENCES cats(id) ON DELETE SET NULL
);

CREATE TABLE targets(
    id UUID PRIMARY KEY,
    mission_id UUID NOT NULL,

    name VARCHAR(255) NOT NULL,
    country VARCHAR(50) NOT NULL,
    notes TEXT DEFAULT '',

    is_completed BOOL NOT NULL DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_mission FOREIGN KEY (mission_id) REFERENCES missions(id) ON DELETE CASCADE
);


CREATE INDEX idx_missions_cat_id ON missions(cat_id);

CREATE INDEX idx_targets_mission_id ON targets(mission_id);

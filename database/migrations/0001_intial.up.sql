CREATE TABLE IF NOT EXISTS imageUpload (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        fileName TEXT NOT NULL,
        url TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        archived_at TIMESTAMP WITH TIME ZONE
)

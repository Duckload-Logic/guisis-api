-- Delete older duplicate educational background records, keeping only the latest (max ID) for each iir_id.
-- This automatically cascades to delete the older school_details.
DELETE eb FROM educational_backgrounds eb
INNER JOIN (
    SELECT iir_id, MAX(id) as max_id
    FROM educational_backgrounds
    GROUP BY iir_id
) keep ON eb.iir_id = keep.iir_id
WHERE eb.id < keep.max_id;

-- Create the unique index on iir_id first (so it can satisfy the foreign key constraint)
CREATE UNIQUE INDEX unique_idx_educational_backgrounds_iir_id ON educational_backgrounds(iir_id);

-- Drop the old non-unique index now that the unique index covers the foreign key constraint
DROP INDEX idx_educational_backgrounds_iir_id ON educational_backgrounds;

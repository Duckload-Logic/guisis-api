-- Restore the non-unique index first (so it can satisfy the foreign key constraint)
CREATE INDEX idx_educational_backgrounds_iir_id ON educational_backgrounds(iir_id ASC);

-- Drop the unique index now that the non-unique index covers the foreign key constraint
DROP INDEX unique_idx_educational_backgrounds_iir_id ON educational_backgrounds;

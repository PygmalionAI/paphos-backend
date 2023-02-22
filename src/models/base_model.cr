abstract class BaseModel < Avram::Model
  def self.database : Avram::Database.class
    AppDatabase
  end

  macro default_columns
    # Override the default primary key type.
    primary_key id : UUID

    # Adds the `created_at` and `updated_at` as usual.
    timestamps
  end
end

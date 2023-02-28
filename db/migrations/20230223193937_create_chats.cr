class CreateChats::V20230223193937 < Avram::Migrator::Migration::V1
  def migrate
    enable_extension "pgcrypto"

    create table_for(Chat) do
      # TODO(11b): Default this to ULID or UUIDv6 to avoid index fragmentation.
      primary_key id : UUID
      add_belongs_to creator : User, on_delete: :restrict # , index: true

      add name : String?
      # add_has_many characters : Character

      add_timestamps
    end
  end

  def rollback
    drop table_for(Chat)
    disable_extension "pgcrypto"
  end
end

class CreateCharacters::V20230223150103 < Avram::Migrator::Migration::V1
  def migrate
    create table_for(Character) do
      primary_key id : Int64
      add slug : String, index: true, unique: true

      add name : String
      add description : String
      add avatar_id : String?

      add greeting : String
      add persona : String
      add world_scenario : String?
      add example_chats : String?

      add visibility : String
      add is_contentious : Bool

      add_belongs_to creator : User, on_delete: :restrict

      add_timestamps
    end
  end

  def rollback
    drop table_for(Character)
  end
end

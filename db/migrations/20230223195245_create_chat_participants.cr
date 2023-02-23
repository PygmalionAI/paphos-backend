class CreateChatParticipants::V20230223195245 < Avram::Migrator::Migration::V1
  def migrate
    create table_for(ChatParticipant) do
      # Needs this useless field until composite primary key support lands in
      # Avram. See:
      # https://github.com/luckyframework/avram/issues/129
      primary_key id : Int64

      add_belongs_to chat : Chat, on_delete: :restrict, foreign_key_type: UUID
      add_belongs_to character : Character, on_delete: :restrict
    end

    create_index :chat_participants, [:chat_id, :character_id], unique: true
  end

  def rollback
    drop table_for(ChatParticipant)
  end
end

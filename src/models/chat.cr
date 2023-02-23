class Chat < BaseModel
  skip_default_columns

  table do
    primary_key id : UUID
    belongs_to creator : User

    column name : String?
    has_many participants : ChatParticipant
    has_many characters : Character, through: [:participants, :character]

    created_at : Time
    updated_at : Time
  end
end

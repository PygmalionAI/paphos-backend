class ChatParticipant < BaseModel
  skip_default_columns

  table do
    # Needs this useless field until composite primary key support lands in
    # Avram. See:
    # https://github.com/luckyframework/avram/issues/129
    primary_key id : Int64

    belongs_to chat : Chat
    belongs_to character : Character
  end
end

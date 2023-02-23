class Character < BaseModel
  VALID_VISIBILITY_VALUES = ["public", "unlisted", "private"]

  table do
    column slug : String

    column name : String
    column description : String
    column avatar_id : String?

    column greeting : String
    column persona : String
    column world_scenario : String?
    column example_chats : String?

    column visibility : String
    column is_contentious : Bool

    belongs_to creator : User
    has_many chat_participations : ChatParticipant
    has_many chats : Chat, through: [:chat_participations, :chat]
  end
end

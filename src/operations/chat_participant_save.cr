class SaveChatParticipant < ChatParticipant::SaveOperation
  permit_columns chat_id, character_id
end

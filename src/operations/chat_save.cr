class SaveChat < Chat::SaveOperation
  needs current_user : User

  param_key :chat

  permit_columns name
  attribute character_slugs : Array(String)

  before_save do
    if name.value != nil
      validate_size_of name, min: 1, max: 32
    end

    assign_creator_id
  end

  after_save do |chat|
    look_up_and_associate_characters(chat)
  end

  private def assign_creator_id
    creator_id.value = current_user.id
  end

  private def look_up_and_associate_characters(chat : Chat)
    return if character_slugs.value == nil

    # ameba:disable Lint/NotNil - guarded by the above return.
    character_slugs.value.not_nil!.each do |slug|
      # TODO(11b): This returns a really undescriptive 404, might be worth
      # refactoring so we can explicitly call out which slug was not found.
      character = CharacterQuery.new.accessible_by(current_user).slug(slug).first
      SaveChatParticipant.create!(chat_id: chat.id, character_id: character.id)
    end
  end
end

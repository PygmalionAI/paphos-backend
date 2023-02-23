class Api::V1::Characters::Update < ApiAction
  include CheckCurrentUser

  patch "/characters/:character_slug" do
    character = CharacterQuery.new.slug(character_slug).first
    ensure_owned_by_current_user!(character)

    # TODO(11b): When chats are implemented, need to block character from being
    # marked "private" if a chat from a user other than the creator already
    # exists.
    #
    # Also, need to block certain fields from being modified I believe.
    updated_character = SaveCharacter.update!(
      character, params, current_user: current_user)

    json({character: FullCharacterSerializer.new(updated_character)})
  end
end

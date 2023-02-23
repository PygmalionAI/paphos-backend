class Api::V1::Characters::Delete < ApiAction
  include CheckCurrentUser

  delete "/characters/:character_slug" do
    character = CharacterQuery.new.slug(character_slug).first
    ensure_owned_by_current_user!(character)

    # TODO(11b): Correctly handle deletion failures due to foreign key
    # constraint failures by sending an adequate error message.
    DeleteCharacter.delete!(character)

    head :ok
  end
end

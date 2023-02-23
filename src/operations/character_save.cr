class SaveCharacter < Character::SaveOperation
  needs current_user : User

  param_key :character

  permit_columns name, description, avatar_id, greeting, persona,
    world_scenario, example_chats, visibility, is_contentious

  before_save do
    validate_size_of name, min: 1, max: 32
    validate_size_of description, min: 8, max: 64
    validate_size_of greeting, min: 2, max: 1024
    validate_size_of persona, min: 12, max: 1024
    validate_size_of world_scenario, max: 512
    validate_size_of example_chats, max: 1024

    validate_inclusion_of visibility, in: Character::VALID_VISIBILITY_VALUES

    # TODO(11b): Consider changing this before going into prod. Apparently this
    # generates an _entire_ UUIDv4 upon collision, which is way overkill.
    Avram::Slugify.set slug, using: name, query: CharacterQuery.new

    creator_id.value = current_user.id
  end
end

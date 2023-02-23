class Api::V1::Characters::Index < ApiAction
  get "/characters" do
    # TODO(11b): pagination, basic filtering (contentious)
    characters = CharacterQuery.new.visible_to(current_user)
    serialized_characters = MinimalCharacterSerializer.for_collection(characters)
    json({characters: serialized_characters})
  end
end

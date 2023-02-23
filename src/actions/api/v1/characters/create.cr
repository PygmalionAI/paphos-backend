class Api::V1::Characters::Create < ApiAction
  post "/characters" do
    character = SaveCharacter.create!(params, current_user: current_user)
    json({character: FullCharacterSerializer.new(character)})
  end
end

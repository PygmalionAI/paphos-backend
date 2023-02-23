# Minimal serializer for the Character model. Meant to be used when the detail
# fields (e.g. required for building a complete prompt) aren't necessary.
class MinimalCharacterSerializer < BaseSerializer
  def initialize(@character : Character)
  end

  def render
    {
      slug:        @character.slug,
      name:        @character.name,
      description: @character.description,
      avatar_id:   @character.avatar_id,
    }
  end
end

class FullCharacterSerializer < BaseSerializer
  def initialize(@character : Character)
  end

  def render
    {
      slug: @character.slug,

      name:        @character.name,
      description: @character.description,
      avatar_id:   @character.avatar_id,

      greeting:       @character.greeting,
      persona:        @character.persona,
      world_scenario: @character.world_scenario,
      example_chats:  @character.example_chats,

      visibility:     @character.visibility,
      is_contentious: @character.is_contentious,
    }
  end
end

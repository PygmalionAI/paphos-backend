class ChatSerializer < BaseSerializer
  def initialize(@chat : Chat)
  end

  def render
    {
      id:         @chat.id,
      name:       @chat.name,
      characters: MinimalCharacterSerializer.for_collection(@chat.characters),
    }
  end
end

class CharacterQuery < Character::BaseQuery
  def visible_to(user : User)
    visibility("public").or(&.creator_id(user.id))
  end
end

class CharacterQuery < Character::BaseQuery
  # Visible characters are characters that are public, or that the user has
  # created.
  def visible_to(user : User)
    where do |where|
      where.visibility("public").or(&.creator_id(user.id))
    end
  end

  # Accessible characters are characters that are public, or unlisted but the
  # user knows about their slug, or private that the user has created
  # themselves.
  def accessible_by(user : User)
    where do |where|
      where.visibility.in(["public", "unlisted"]).or(&.creator_id(user.id))
    end
  end
end

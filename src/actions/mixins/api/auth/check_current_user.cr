module CheckCurrentUser
  class UnauthorizedError < Lucky::Error
  end

  # TODO(11b): Create an admin type that can bypass these checks.

  def ensure_owned_by_current_user!(character : Character)
    if character.creator_id != current_user.id
      raise UnauthorizedError.new
    end
  end

  def ensure_owned_by_current_user!(chat : Chat)
    if chat.creator_id != current_user.id
      raise UnauthorizedError.new
    end
  end
end

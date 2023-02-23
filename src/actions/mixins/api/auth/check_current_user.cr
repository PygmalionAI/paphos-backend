module CheckCurrentUser
  class UnauthorizedError < Lucky::Error
  end

  def ensure_owned_by_current_user!(character : Character)
    # TODO(11b): Create an admin type that can bypass this check.
    if character.creator_id != current_user.id
      raise UnauthorizedError.new
    end
  end
end

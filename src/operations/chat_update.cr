class UpdateChat < Chat::SaveOperation
  needs current_user : User

  param_key :chat

  permit_columns name

  before_save do
    # TODO(11b): How do I stop `null` from being converted into a String at the
    # controller level, so I don't have to do this?
    if name.value == "null"
      name.value = nil
    end

    if name.value != nil
      validate_size_of name, min: 1, max: 32
    end
  end
end

class Api::V1::Users::Login::Create < ApiAction
  include Api::Auth::SkipRequireAuthToken

  post "/users/login" do
    LogInUser.run(params) do |operation, user|
      if user
        json({token: UserToken.generate(user)})
      else
        raise Avram::InvalidOperationError.new(operation)
      end
    end
  end
end

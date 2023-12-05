# Homework Screenshot
## Edit - PUT /users/:id
- 个人昵称长度超出20限制
  ![img.png](postman-shots/week-02/edit-nickname-exceed-limitation.png)
- 生日格式不正确
  ![img_1.png](postman-shots/week-02/edit-birthdate-incorrect-01.png)
  ![img.png](postman-shots/week-02/edit-birthdate-incorrect-02.png)
- 个人简介长度超出150限制
  ![img.png](postman-shots/week-02/edit-personal-profile-exceed-limitation.png)
- 成功编辑后返回完整的个人信息
  ![img.png](postman-shots/week-02/edit-succeed-return-user-info.png)
  ![img.png](postman-shots/week-02/edit-succeed-db.png)
## Profile - GET /users/:id
- 通过存在的ID获取用户profile
  ![img.png](postman-shots/week-02/profile-succeed-return-user-info.png)
- 通过不存在的ID获取用户profile
  ![img.png](postman-shots/week-02/profile-error-return-not-found.png)

describe('CommWithGql Component', () => {
  beforeEach(() => {
    cy.visit('http://localhost:8888/#/manager');
    // cy.wait(2000); // 等待2秒以確保數據加載完成
  });

  // it('should display a list of items', () => {
  //   cy.get('.listing-item').should('have.length.greaterThan', 0);
  // });

  it('should click the button and check if the menu is correct', () => {
    cy.contains('button', '新建會議室').should('exist').click();;
    cy.contains('button', '編輯').should('exist');
    cy.contains('button', '取消').should('exist');
    cy.contains('label', '會議名稱：')
      .parent()
      .find('input')
      .should('exist')
      .clear()
      .type("meeting Test Room 1");;
    cy.contains('label', '圖片網址')
      .parent()
      .find('input')
      .should('exist');
    cy.contains('label', '人數限制')
      .parent()
      .find('input')
      .should('exist')
      .clear()
      .type("5");
    cy.contains('label', '可否進食：')
      .parent()
      .find('input[type="checkbox"]')
      .should('exist');
    cy.contains('label', '有大桌子：')
      .parent()
      .find('input[type="checkbox"]')
      .should('exist');
    cy.contains('label', '有投影機：')
      .parent()
      .find('input[type="checkbox"]')
      .should('exist');
    cy.contains('button', '編輯').click();
    cy.contains('meeting Test Room 1').should('exist');
    cy.get('.main-white-button').find('a').contains('編輯');
    cy.get('.main-white-button').find('a').contains('刪除').click();

  });

});

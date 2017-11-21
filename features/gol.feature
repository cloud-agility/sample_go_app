# file: $GOPATH/godogs/features/godogs.feature
Feature: game of life
  In order to play game of life
  As board
  I need to evolve an input state to produce an evolved state

  Scenario: empty board stays empty
    Given an empty board
    When I evolve it
    Then it should be empty

  Scenario: single cell board dies
    Given a single cell on the board
    When I evolve it
    Then it should be empty

  Scenario: 3x3 board with single cell dies
    Given a 3 x 3 board with the following
    """
    ...
    .*.
    ...
    """
    When I evolve it
    Then it should be empty

  Scenario: 3x3 board with stable block
    Given a 3 x 3 board with the following
    """
    ...
    .**
    .**
    """
    When I evolve it
    Then it should be like the following
    """
    ...
    .**
    .**
    """
  Scenario: 6x6 board with nascent stable block
    Given a 6 x 6 board with the following
    """
    ......
    ......
    ..**..
    ..*...
    ......
    ......
    """
    When I evolve it
    Then it should be like the following
    """
    ......
    ......
    ..**..
    ..**..
    ......
    ......
    """
